package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	types "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	//TODO:connect to k8s by client-go
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, "config", "guyu"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
	/*TODO:if you cannot connect to locol k8s , Notes the panic part ,then this part can test the sort
	h, _ := time.ParseDuration("-1h")
	pods = &types.PodList{
		Items: []types.Pod{
			types.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "pod1",
					Namespace:         "ns1",
					CreationTimestamp: metav1.Time{time.Now().Add(2 * h)},
				},
			},
			types.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "pod2",
					Namespace:         "ns2",
					CreationTimestamp: metav1.Time{time.Now()},
				},
			},
			types.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "pod3",
					Namespace:         "ns3",
					CreationTimestamp: metav1.Time{time.Now().Add(h)},
				},
			},
		},
	}
	*/
	//TODO:pod根据age排序
	podRatios := newPodRatios(pods)
	sort.Sort(sort.Reverse(podRatios))

	n := len(podRatios) / 2
	fmt.Printf("Namespace %s PodName %s Age %d N %d ", podRatios[n].NameSpace, podRatios[n].PodName, podRatios[n].Age, n)
}
func homeDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

type podRatio struct {
	PodName   string
	NameSpace string
	Age       float64
}

//实现golang sort包，内涵插入排序．快排和堆排序三种排序算法

type podRatios []podRatio

func (pr podRatios) Len() int {
	return len(pr)
}
func (pr podRatios) Less(i, j int) bool {
	return pr[i].Age < pr[j].Age
}
func (pr podRatios) Swap(i, j int) {
	pr[i], pr[j] = pr[j], pr[i]
}

func newPodRatios(podList *types.PodList) podRatios {
	pr := make(podRatios, 0, len(podList.Items))
	for _, pod := range podList.Items {
		pr = append(pr, podRatio{
			PodName:   pod.Name,
			NameSpace: pod.Namespace,
			Age:       time.Now().Sub(pod.CreationTimestamp.Time).Seconds(),
		})
	}
	fmt.Printf("guyutest %v \n", pr)
	return pr
}
