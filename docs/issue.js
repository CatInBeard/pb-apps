async function getIssues(owner, repo, cacheTimeMinutes = 30) {
    const cacheKey = `github_issues_${owner}_${repo}`;
    const cachedData = localStorage.getItem(cacheKey);

    if (cachedData) {
        try {
            const { timestamp, data } = JSON.parse(cachedData);
            const now = Date.now();
            const cacheExpirationTime = timestamp + cacheTimeMinutes * 60 * 1000;

            if (now < cacheExpirationTime) {
                console.log("Данные о issues взяты из кэша");
                displayIssues(data);
                return;
            } else {
                localStorage.removeItem(cacheKey);
                console.log("Кэш issues устарел, удаляем");
            }
        } catch (error) {
            console.error("Ошибка при разборе данных о issues из кэша:", error);
            localStorage.removeItem(cacheKey);
        }
    }

    try {
        const response = await fetch(`https://api.github.com/repos/${owner}/${repo}/issues?state=all`);

        if (!response.ok) {
            throw new Error(`Failed to use github api: ${response.status} ${response.statusText}`);
        }

        const issues = await response.json();

        // Кэшируем данные о issues
        const dataToCache = {
            timestamp: Date.now(),
            data: issues
        };

        localStorage.setItem(cacheKey, JSON.stringify(dataToCache));
        console.log("Данные о issues сохранены в кэш");

        displayIssues(issues);
    } catch (error) {
        console.error(error);
    }
}

async function getComments(url, container) {
    const cacheKey = `github_comments_${url}`;
    const cachedData = localStorage.getItem(cacheKey);

    if (cachedData) {
        try {
            const { timestamp, data } = JSON.parse(cachedData);
            const now = Date.now();
            const cacheExpirationTime = timestamp + 30 * 60 * 1000; // Кэшируем комментарии на 30 минут

            if (now < cacheExpirationTime) {
                console.log("Комментарии взяты из кэша");
                displayComments(data, container);
                return;
            } else {
                localStorage.removeItem(cacheKey);
                console.log("Кэш комментариев устарел, удаляем");
            }
        } catch (error) {
            console.error("Ошибка при разборе данных о комментариях из кэша:", error);
            localStorage.removeItem(cacheKey);
        }
    }

    try {
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`Failed to fetch comments: ${response.status} ${response.statusText}`);
        }
        const comments = await response.json();

        // Кэшируем данные о комментариях
        const dataToCache = {
            timestamp: Date.now(),
            data: comments
        };

        localStorage.setItem(cacheKey, JSON.stringify(dataToCache));
        console.log("Комментарии сохранены в кэш");

        displayComments(comments, container);
    } catch (error) {
        console.error(error);
    }
}

function displayComments(comments, container) {
    comments.forEach(comment => {
        const commentElement = document.createElement('div');
        commentElement.className = 'comment';
        commentElement.innerHTML = `
            <p>${comment.body}</p>
            <p>Commented by ${comment.user.login} on ${new Date(comment.created_at).toLocaleString()}</p>
        `;
        container.appendChild(commentElement);
    });
}


function displayIssues(issues) {
    const issuesContainer = document.getElementById('issues-container');
    const issuesHeader = document.getElementById('issues-header');

    let hasIssues = false;
    issuesContainer.innerHTML = '';

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
        if (issuesHeader) {
            issuesHeader.classList.remove("d-none");
        }
    } else {
        issuesContainer.innerHTML = '<p>No issues found.</p>';
    }
}


window.addEventListener('load', () => {
    getIssues("CatInBeard", "pb-apps");
});
