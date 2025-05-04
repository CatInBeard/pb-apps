import os
import requests
import json
import sys

github_token = os.environ.get('GITHUB_TOKEN')

if github_token is None:
    sys.stderr.write("Token not found\n")

repo_owner = 'CatInBeard'
repo_name = 'pb-apps'

issues_url = f'https://api.github.com/repos/{repo_owner}/{repo_name}/issues'

headers = {}
if github_token:
    headers['Authorization'] = f'token {github_token}'

response = requests.get(issues_url, headers=headers, params={'state': 'all', 'per_page': 100})

if response.status_code == 200:
    issues_data = response.json()
    issues = []
    for issue in issues_data:
        if 'pull_request' not in issue:
            issue_comments_url = f'https://api.github.com/repos/{repo_owner}/{repo_name}/issues/{issue["number"]}/comments'
            response_comments = requests.get(issue_comments_url, headers=headers, params={'per_page': 100})
            if response_comments.status_code == 200:
                comments_data = response_comments.json()
                issue['comments'] = comments_data
            issues.append(issue)
    issues.sort(key=lambda x: x['created_at'])
    with open('issues.json', 'w') as f:
        json.dump(issues, f, indent=4)
    print("Successfully fetched issues")
else:
    sys.stderr.write("Can't get issues!\n")
    exit(1)
