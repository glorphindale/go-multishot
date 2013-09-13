#!/usr/bin/python

import unittest
import requests
import collections

Response = collections.namedtuple("Response", ['code', 'body'])

TARGET="http://localhost:8080"

def parse_reponse(raw_response):
    body = {}
    if raw_response.status_code == 200:
        body = {k: v for k, v in [s.split(":", 1) for s in raw_response.text.split("\n") ]}
    return Response(raw_response.status_code, body)

class SimpleTests(unittest.TestCase):
    def test_get(self):
        r = requests.get(TARGET + "/get")
        response = parse_reponse(r)

        self.assertEqual(response.code, 200)
        self.assertEqual(response.body["Request"], "/get")
        self.assertEqual(response.body["Body"], "<Empty>")

    def test_user_agent(self):
        r = requests.get(TARGET + "/get-agent", headers={"User-Agent": "My User Agent"})
        response = parse_reponse(r)

        self.assertEqual(response.code, 200)
        self.assertTrue("My User Agent" in response.body["Headers"], "Wrong user agent in '" +
                response.body["Headers"] + "'")

    def test_get_query(self):
        r = requests.post(TARGET + "/get-query?param1=val1&param2=val2")
        response = parse_reponse(r)

        self.assertEqual(response.code, 200)
        self.assertEqual(response.body["Request"], "/get-query?param1=val1&param2=val2")
        self.assertEqual(response.body["Body"], "<Empty>")

    def test_post(self):
        r = requests.post(TARGET + "/post", data="POST body")
        response = parse_reponse(r)

        self.assertEqual(response.code, 200)
        self.assertEqual(response.body["Request"], "/post")
        self.assertEqual(response.body["Body"], "POST body")


if __name__ == "__main__":
    unittest.main()
