OH
==

Why
====
Integrate Airflow with Openshift.
Problem: Desirable to have separate pods for webserver and scheduler of Airflow, but no shared storage available.
Needed replecation of DAGs between two services.

From other side I want to place all replication logic in one pod. That's why Go.

Planned replication functionality:
1. Upload DAGs
2. Monitor DAGs modifications in web dashboard.

*Oh - a hero of Ukrainian fairy tail
