# Metronomikon

## Overview

Metronomikon is a Metronome API compatibility layer for Kubernetes. It responds to the Metronome V2 API and maps
requests into equivalent Kubernetes API calls.

### Why would anyone want this?

Our motivation for writing this was to ease our transition from DC/OS to Kubernetes. Many of our applications
are designed to interact with the Metronome API in DC/OS, and it was easier to transition these services to
Kubernetes if they didn't need to be modified and could work as-is in both DC/OS and Kubernetes.

### What about the differences between the way that Metronome and Kubernetes CronJobs work?

While it's true that Metronome and Kubernetes CronJobs work in fairly different ways, Metronomikon provides
various methods to help work around these differences. In places where the Kubernetes API provides a richer
variety of knobs and dials, Metronomikon has configurable defaults for these various values where needed. In
cases where Metronome provides greater functionality (such as the ability to trigger a job "now"), the feature
will be approximated in Kubernetes (by creating a Job based on a CronJob template in the case of triggering a
job "now").

Unfortunately, it won't be possible to expose the Kubernetes only options through the Metronome API for
compatibility reasons. It's also possible that some of the Metronome functionality may not be available due
to an inability to reasonable approximate that behavior in Kubernetes.

One such unsupported Metronome feature, is job "stop" functionality. Kubernetes CronJobs does not support stopping a running CronJob without also deleting the job instance itself, which would in turn destroy history of said job.
