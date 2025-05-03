
async function getIssues(owner, repo) {
    try {
        const response = await fetch(`https://api.github.com/repos/${owner}/${repo}/issues?state=all`);

        if (!response.ok) {
            throw new Error(`Failed to use github api: ${response.status} ${response.statusText}`);
        }

        const issues = await response.json();

        const issuesContainer = document.getElementById('issues-container');

        let hasIssues = false;

        issues.forEach(issue => {
            if (issue.pull_request === undefined) {
                hasIssues = true;

                const issueElement = document.createElement('div');
                issueElement.className = 'issue';

                const issueHeader = document.createElement('div');
                issueHeader.className = 'issue-header';
                issueHeader.innerHTML = `
                                <h4>${issue.title}</h4>
                                <p>Opened by ${issue.user.login} on ${new Date(issue.created_at).toLocaleString()}</p>
                            `;
                issueElement.appendChild(issueHeader);

                const issueComments = document.createElement('div');
                issueComments.className = 'issue-comments';

                const issueBody = document.createElement('div');
                issueBody.className = 'issue-body';
                issueBody.innerHTML = `
                                <p>${issue.body}</p>
                            `;
                issueComments.appendChild(issueBody);

                getComments(issue.comments_url, issueComments);

                issueElement.appendChild(issueComments);

                issuesContainer.appendChild(issueElement);
            }
        });

        if (hasIssues) {
            issuesContainer.classList.remove('d-none');
            issuesHeader.remove("d-none");
        } else {
            issuesContainer.innerHTML = '<p>No issues found.</p>';
        }
    } catch (error) {
        console.error(error);
    }
}

async function getComments(url, container) {
    try {
        const response = await fetch(url);
        const comments = await response.json();

        comments.forEach(comment => {
            const commentElement = document.createElement('div');
            commentElement.className = 'comment';
            commentElement.innerHTML = `
                            <p>${comment.body}</p>
                            <p>Commented by ${comment.user.login} on ${new Date(comment.created_at).toLocaleString()}</p>
                        `;
            container.appendChild(commentElement);
        });
    } catch (error) {
        console.error(error);
    }
}

window.addEventListener('load', () => {
    getIssues("CatInBeard", "pb-apps");
});
