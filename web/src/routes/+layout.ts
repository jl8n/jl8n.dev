export type NowPlaying = {
    Artist: string;
    Album: string;
    Track: string;
    AlbumMBID: string;
    Art: string;
    ArtAvailable: boolean;
};

export async function load({ fetch }) {
    const res = await fetch('http://localhost:3000/');

    
    if (!res.ok) {
        throw new Error(`An error occurred: ${res.status}`);
    }

    const data = await res.json() as NowPlaying;
    console.log('data', data)
    let blobUrl = ''

    if (data.ArtAvailable) {
        const params = `filename=${encodeURIComponent(data.AlbumMBID)}`
        const res2 = await fetch('http://localhost:3000/albumart?' + params);
        const blob = await res2.blob()
        blobUrl = URL.createObjectURL(blob);
    }

    return {
        Artist: data.Artist,
        Album: data.Album,
        Track: data.Track,
        Art: blobUrl
    };
}