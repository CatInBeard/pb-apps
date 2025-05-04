import os
import requests
import json
import sys

github_token = os.environ.get('GITHUB_TOKEN')

if github_token is None:
        sys.stderr.write("Token not found\n")

repo_json_url = 'https://raw.githubusercontent.com/CatInBeard/pb-apps/refs/heads/main/repo.json'

response = requests.get(repo_json_url)
if response.status_code == 200:
    repo_data = response.json()
else:
    sys.stderr.write("Can't get repo.json!\n")
    exit(1)

for repo in repo_data['repositories']:
    repo_name = repo['name']
    repo_owner = repo['url'].split('/')[-2]

    release_url = f'https://api.github.com/repos/{repo_owner}/{repo_name}/releases/latest'

    headers = {}
    if github_token:
        headers['Authorization'] = f'token {github_token}'
    response = requests.get(release_url, headers=headers)

    if response.status_code == 200:
        release_data = response.json()
        repo['version'] = release_data['tag_name']
        repo['release_link'] = release_data['zipball_url']
    else:
        sys.stderr.write("Repo " + repo_name + "/" + repo_owner+ " returns " + str(response.status_code) +" \n")
        exit(1)
    
with open('repo.json', 'w') as f:
    json.dump(repo_data, f, indent=4)
print("Successfully update repo")