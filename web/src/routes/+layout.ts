export type NowPlaying = {
    Artist: string;
    Album: string;
    Track: string;
};

export async function load({ fetch }) {
    const res = await fetch('http://localhost:3000/');
    const data = await res.json() as NowPlaying;

    return {
        Artist: data.Artist,
        Album: data.Album,
        Track: data.Track
    };
}