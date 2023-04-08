import { NextApiRequest, NextApiResponse } from "next"

const API_KEY = process.env.TMDB_API_KEY

interface Poster {
    "aspect_ratio": number,
    "height": number,
    "iso_639_1": string,
    "file_path": string,
    "vote_average": number,
    "vote_count": number,
    "width": number
}

function getBestPosterPath(posters: Poster[]) {
    const poster = posters.find(poster => poster.iso_639_1 === "en")

    if (!poster) {
        return posters[0].file_path
    }

    return poster.file_path
}

export default async function handler(
    req: NextApiRequest,
    res: NextApiResponse
    ) {
    const { id, size } = req.query
    const response = await fetch(
        `https://api.themoviedb.org/3/movie/${id}/images?api_key=${API_KEY}`
    )
    const data = await response.json()
    const posterPath = getBestPosterPath(data.posters)
    const posterUrl = `https://image.tmdb.org/t/p/w${size}${posterPath}`
    res.redirect(posterUrl)
}