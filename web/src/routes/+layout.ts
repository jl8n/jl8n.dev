// export const prerender = true;  // adapter-static
export const prerender = true

// export type NowPlaying = {
//     Artist?: string;
//     Album?: string;
//     Track?: string;
//     AlbumMBID?: string;
//     Art?: string;
//     ArtAvailable?: boolean | null;
// };

// export async function load({ fetch }) {
//     let res;
//     try {
//         res = await fetch('http://localhost:3000/');
//     } catch (error) {
//         console.error(`Fetch failed: ${error}`);
//         return;
//     }

//     if (!res.ok) {
//         console.error(`An error occurred: ${res.status}`);
//         return;
//     } else if (res.status == 204) {
//         return {} as NowPlaying
//     }

//     let data: NowPlaying = {};

//     try {
//         const resData = await res.json() as NowPlaying;
//         data = {
//             Artist: resData.Artist,
//             Album: resData.Album,
//             Track: resData.Track,
//             AlbumMBID: resData.AlbumMBID,
//             ArtAvailable: resData.ArtAvailable
//         }
//     } catch (e) {
//         console.error(e);
//     }

//     if (data.ArtAvailable && data.AlbumMBID) {
//         const params = `filename=${encodeURIComponent(data.AlbumMBID)}`;
//         let res2;
//         try {
//             res2 = await fetch('http://localhost:3000/albumart?' + params);
//         } catch (error) {
//             console.error(`Fetch failed: ${error}`);
//             return;
//         }
        
//         const blob = await res2.blob();
//         data.Art = URL.createObjectURL(blob);
//     }

//     console.log('final data', data);

//     return data;
// }
