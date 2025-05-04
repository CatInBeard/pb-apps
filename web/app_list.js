const buildAppList = () => {
    const appListContainer = document.getElementById('app-list');
    const appListBlock = document.getElementById('app-list-block');

    fetch('https://raw.githubusercontent.com/CatInBeard/pb-apps/main/repo.json')
        .then(response => response.json())
        .then(data => {
            const apps = data.repositories;

            if (apps && apps.length > 0) {
                appListContainer.classList.remove('d-none');
                appListBlock.classList.remove('d-none');
                apps.forEach(async app => {
                    if (app.name == "pb-apps") {
                        return;
                    }
                    const appItem = document.createElement('div');
                    appItem.classList.add('card');

                    const appItemBody = document.createElement('div');
                    appItemBody.classList.add('card-body');

                    const nameElement = document.createElement('h5');
                    nameElement.textContent = app.name;

                    const descriptionElement = document.createElement('p');
                    descriptionElement.textContent = app.description;

                    const urlElement = document.createElement('a');

                    let url = await getRedirectTarget(app.url + "/releases/latest")

                    if(url) {
                        url = url + "/release.zip"
                    } else {
                        url = app.url + "/releases/latest";
                    }

                    urlElement.href = url;
                    urlElement.textContent = 'Download';
                    urlElement.classList.add('btn', 'btn-primary', 'btn-sm');
                    urlElement.target = "_blank";

                    const licenseElement = document.createElement('p');
                    licenseElement.textContent = `License: ${app.license}`;
                    licenseElement.classList.add('text-muted');

                    appItemBody.appendChild(nameElement);
                    appItemBody.appendChild(descriptionElement);
                    appItemBody.appendChild(urlElement);
                    appItemBody.appendChild(licenseElement);

                    appItem.appendChild(appItemBody)

                    appListContainer.appendChild(appItem);
                });
            }
        })
        .catch(error => {
            console.error('Error fetching data:', error);
        });
}

async function getRedirectTarget(url) {
    try {
        const response = await fetch(url, {
            method: 'HEAD',
            redirect: 'manual'
        });

        if (response.status >= 300 && response.status < 400) {
            const redirectUrl = response.headers.get('Location');
            if (redirectUrl) {
                return redirectUrl;
            } else {
                return null;
            }
        } else {
            return null;
        }
    } catch (error) {
        console.error(error);
        return null;
    }
}

window.addEventListener('load', () => {
    buildAppList()
});