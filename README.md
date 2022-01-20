# SpiderSwarm

SpiderSwarm project is developing a distributed cloud-native system for scalable web scraping. Based on dataflow programming concepts, it aims to
enable declaratively defining web scraping workflows and distributing them across a large number of worker nodes in a frictionless way. 
SpiderSwarm is meant to be scalable from a single node running on low-powered systems to thousands of nodes in the cloud deployment. 
Furthermore, it exposes APIs and can be integrated into end-user-facing apps.

There are the following key actor subsystems to SpiderSwarm:

* Manager - schedules tasks based on workflow and received task results. If task results contain items(s) (pieces of scraped data) it is sent to Exporter.
* Worker - received scheduled tasks from Manager, executed them, and sends back task results.
* Exporter - receives items and writes them out to the external store. At this point, only CSV files are supported.

The workflow consists of one or more task templates. Each task template is consisting of one or more action templates, connected by data pipe templates.
When the Manager starts executing a Workflow, it schedules a single task from the template that has been marked as initial. When Worker received a
scheduled task object it fleshes out the task by instantiating actions and data pipes (queues for data). Once all actions are executed, the task result
is created. It may contain some items and some promises for tasks to be launched later. The Manager receives the task result and unwraps it. If there
are any items, they are sent to Exporter. If there are any task promises, the corresponding new tasks (based on templates in the workflow) are being
scheduled. This process is being repeated until there is no further work to do.

At this point, we use Redis for task queue/message bus, but support for more underlying technologies for messaging (Zookeeper? Kafka?) is planned.

Run the following command to build a Docker image:
```
docker build -t spiderswarm:0.0.0 .
```

