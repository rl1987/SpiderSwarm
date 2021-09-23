Glossary
========

* **Action** - a discrete unit of work that is performed as part of **Task**.
* **Exporter** - a component of SpiderSwarm that exports scraped data to external data store (CSV file, spreadsheet, database, etc.).
* **DataChunk** - an object that holds one or more pieces of intermediate information.
* **DataPipe** - a queue object that stores intermediate data between **Action**s.
* **Task** - a building block for **Workflow** that is being executed by **Worker** and results in zero or more **Item**s and/or zero or more new **Task**s.
* **TaskPromise** - a notice about task to be done in the future.
* **Workflow** - scraping program that is executed by SpiderSwarm distributed system.
* **Item** - a scraped datapoint that consists of one or more key-value pairs.
* **ScheduledTask** - an object that encapsulates **TaskPromise** and **TaskTemplate** and is sent to **Worker**.
* **SpiderBus** - a message bus and task queue for data transfer and workload distribution.
* **SpiderBusBackend** - an interface between **SpiderBus** and underlying system (database, message queue) that **SpiderBus** uses as foundation.

