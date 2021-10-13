from behave import *
from behave.contrib.scenario_autoretry import patch_scenario_with_autoretry


def before_feature(context, feature):
    '''
    Automatically wrap all scenarios tagged with 'autoretry' with a retry loop
    '''
    for scenario in feature.walk_scenarios():
        if "autoretry" in scenario.effective_tags:
            patch_scenario_with_autoretry(scenario, max_attempts=60)
