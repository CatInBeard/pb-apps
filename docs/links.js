async function getLatestReleaseZipUrl(owner, repo, cacheTimeMinutes = 90) {
    const cacheKey = `github_release_${owner}_${repo}`;
    const cachedData = localStorage.getItem(cacheKey);

    if (cachedData) {
        const { timestamp, data } = JSON.parse(cachedData);
        const now = Date.now();
        const cacheExpirationTime = timestamp + cacheTimeMinutes * 60 * 1000;

        if (now < cacheExpirationTime) {
            return data;
        } else {
            localStorage.removeItem(cacheKey);
        }
    }

    try {
        const response = await fetch(`https://api.github.com/repos/${owner}/${repo}/releases/latest`);

        if (!response.ok) {
            throw new Error(`Failed to use github api: ${response.status} ${response.statusText}`);
        }

        const release = await response.json();

        const zipAsset = release.assets.find(asset => asset.name.endsWith('.zip'));

        if (!zipAsset) {
            throw new Error('release.zip not found');
        }

        const zipUrl = zipAsset.browser_download_url;

        const dataToCache = {
            timestamp: Date.now(),
            data: zipUrl,
        };
        localStorage.setItem(cacheKey, JSON.stringify(dataToCache));

        return zipUrl;

    } catch (error) {
        console.error(error);
        return null;
    }
}


window.addEventListener('load', () =>  {
    getLatestReleaseZipUrl("CatInBeard", "pb-apps")
        .then(zipUrl => {
            if (zipUrl) {
                release_link.href = zipUrl
            }
        });
});