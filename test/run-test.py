#!/usr/bin/env python

from __future__ import print_function

import json
import requests
import sys
import time

def compare_structures(data1, data2):
    if type(data1) != type(data2):
        return False
    if isinstance(data1, list):
        if len(data1) != len(data2):
            return False
        for idx in range(0, len(data1)):
            if not compare_structures(data1[idx], data2[idx]):
                return False
    elif isinstance(data2, dict):
        if sorted(data1.keys()) != sorted(data2.keys()):
            return False
        for key in data1:
            if not compare_structures(data1[key], data2[key]):
                return False
    else:
        return (data1 == data2)
    return True

def perform_request(test_data, url):
    if test_data['method'] == 'GET':
        r = requests.get(url)
    elif test_data['method'] == 'PUT':
        r = requests.put(url)
    elif test_data['method'] == 'POST':
        r = requests.post(url)
    elif test_data['method'] == 'DELETE':
        r = requests.delete(url)
    else:
        print('Unsupported HTTP method: %s' % test_data['method'])
        return False

    if r.status_code != test_data['responseCode']:
        print('Got response code %d, expected %d' % (r.status_code, test_data['responseCode']))
        return False

    if 'responseText' in test_data:
        if r.text != test_data['responseText']:
            print('Got response text:\n\n%s\n\nexpected:\n\n%s' % (r.text, test_data['responseText']))
            return False
    elif 'responseJsonFile' in test_data:
        response_data = json.load(open('%s/%s' % (TEST_DIR, test_data['responseJsonFile'])))
        if not compare_structures(r.json(), response_data):
            print('Got response JSON:\n\n%s\n\nexpected:\n\n%s' % (r.text, json.dumps(response_data)))
            return False
    return True

ENDPOINT = sys.argv[1]
TEST_DIR = sys.argv[2]

test_data = json.load(open('%s/metadata.json' % TEST_DIR))

retries = test_data.get('retries', 1)

print('Running test step "%s"...' % test_data['name'])

url = '%s%s' % (ENDPOINT, test_data['urlPath'])

print('Performing %s request against URL %s' % (test_data['method'], url))

for i in range(retries):
    if perform_request(test_data, url):
        sys.exit(0)
    print("Try %s in %s failed" % (i+1, retries))
    time.sleep(1)
sys.exit(1)
