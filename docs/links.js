async function getLatestReleaseZipUrl(owner, repo) {
    try {
        const response = await fetch(`https://api.github.com/repos/${owner}/${repo}/releases/latest`);

        if (!response.ok) {
            throw new Error(`Failed to use github api: ${response.status} ${response.statusText}`);
        }

        const release = await response.json();

        const zipAsset = release.assets.find(asset => asset.name.endsWith('.zip'));

        if (!zipAsset) {
            console.log('release.zip not found');
            return null;
        }

        const zipUrl = zipAsset.browser_download_url;
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