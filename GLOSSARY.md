Glossary
========

* **Action** - a discrete unit of work that is performed as part of **Task**.
* **Exporter** - a component of SpiderSwarm that exports scraped data to external data store (CSV file, spreadsheet, database, etc.).
* **DataPipe** - a queue object that stores intermediate data between **Action**s.
* **Task** - a building block for **Workflow** that is being executed by **Worker** and results in zero or more **Item**s and/or zero or more new **Task**s.
* **Workflow** - scraping program that is executed by SpiderSwarm distributed system.
* **Item** - a scraped datapoint that consists of one or more key-value pairs.

