Feature: Job management

  Scenario: Get all jobs
    When we request all jobs
    Then we should see the expected jobs returned
      | namespace          | name   | command                                                                       |
      | metronomikon-test1 | hello1 |                                                                               |
      | metronomikon-test1 | hello2 |                                                                               |
      | metronomikon-test2 | hello3 |                                                                               |
      | metronomikon-test2 | hello4 | date; echo Hello from the Kubernetes cluster hello4 - I will now fail; exit 1 |
      | metronomikon-test3 | hello5 |                                                                               |

  Scenario: Get job
    When we request a specific job
      | job                       |
      | metronomikon-test1.hello2 |
    Then we should see the expected job returned
      | namespace          | name   |
      | metronomikon-test1 | hello2 |

  Scenario: Get job with non-standard command
    When we request a specific job
      | job                       |
      | metronomikon-test2.hello4 |
    Then we should see the expected job returned
      | namespace          | name   | command                                                                       |
      | metronomikon-test2 | hello4 | date; echo Hello from the Kubernetes cluster hello4 - I will now fail; exit 1 |

  Scenario: Delete job
    When we delete a specific job
      | job                       |
      | metronomikon-test3.hello5 |
    Then we should see the expected job returned
      | namespace          | name   |
      | metronomikon-test3 | hello5 |
     And we expect the job to be deleted

  Scenario: Check that non-existing job fails
    When we delete a specific job
      | job                     |
      | notexisting.notexisting |
    Then we expect the job to not be found

  Scenario: Run a job
    When we run a specific job
      | job                       |
      | metronomikon-test1.hello1 |
    Then we check that the job has completed
     And we check that the completed run shows up in the list of job runs

  Scenario: Run a job that should fail
    When we run a specific job
      | job                       |
      | metronomikon-test2.hello4 |
    Then we check that the job has failed

  Scenario: Get job history
    When we fetch history for a specific job
      | job                       |
      | metronomikon-test1.hello1 |
    Then we check that the job history has the expected output

  @autoretry
  Scenario: Get history for all jobs
    When we fetch history for all jobs
    Then we check that the history for all jobs has the expected output
      | namespace          | name   | failed | command                                                                       |
      | metronomikon-test1 | hello1 |        |                                                                               |
      | metronomikon-test1 | hello2 |        |                                                                               |
      | metronomikon-test2 | hello3 |        |                                                                               |
      | metronomikon-test2 | hello4 |  true  | date; echo Hello from the Kubernetes cluster hello4 - I will now fail; exit 1 |
