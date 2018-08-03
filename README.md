OH
====

## Why

Integrate Airflow with Openshift.
Problem: Desirable to have separate pods for webserver and scheduler of Airflow, but no shared storage available.
Needed replecation of DAGs between two services.

From other side I want to place all replication logic in one pod. That's why Go.

## Planned replication functionality:

  - Upload DAGs
  - Monitor DAGs modifications in web dashboard.

### Project name

![Oh image](http://book4u.in.ua/upload/books/m/wh/855f3fe73f32.jpg)
Oh - a hero of Ukrainian fairy tail.
