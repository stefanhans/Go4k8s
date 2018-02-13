# Go4k8s

The [Showcase directory](./Showcase) refers to the talk "[K8s API & Go Programming](http://go-talks.appspot.com/github.com/stefanhans/go-present/slides/Kubernetes/IntroductionIntoClient-Go.slide#1)".

The [Examples directory](./Examples) serves as playground with the official examples from [kubernetes/client-go](https://github.com/kubernetes/client-go).

### Lessons Learned

- FYI: To find out what kubectl is doing under the hood you can use "--v=10" for instance. Then, all calls of glog.V(level) are displayed together with some more information.

- FYI: The Go tool godoc is a nice CLI for GoDoc and code comments, respectively. It has a html version for the browser and interesting outputs in the terminal.

  godoc -http=:6060 -index # let you search at localhost:6060

  Probably better is it with index file:

  godoc -index_files index.file -write_index # one time preparation

  godoc -http=:6060 -index -index_files index.file
  or
  godoc -index -index_files index.file -q ServiceInterface

