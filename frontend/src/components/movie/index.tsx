import { AiFillStar } from "react-icons/ai"

import Image from 'next/image'

export interface Movie {
    Title: string
    TMBdID: number
    ReleaseDate: string
    RatingAverage: number
    RatingCount: number
    ReleaseStatus: string
    Overview: string
}

interface MovieProps {
    movie: Movie
    posterPath?: string
    size?: string
    id: number
}

export function MovieItem({ movie }: MovieProps) {

 const formatDate = (date: string) => {
    const dateObj = new Date(date)
    const month = dateObj.toLocaleString('default', { month: 'long' })
    const day = dateObj.getDate()
    const year = dateObj.getFullYear()

    return `${month} ${day}, ${year}`
}

const posterUrl = `http://localhost:3000/api/movie_image?id=${movie.TMBdID}&size=${500}`
  
  return (
    <li className="flex rounded overflow-hidden border hover:bg-blue-100">
        <div className="max-w-md">
            <Image width={64} height={64} src={posterUrl} alt={movie.Title}/>
        </div>
        <div className="flex flex-col justify-between p-2 w-full">
            <div className="flex flex-col"> 
                <span className="text-slate-500">{movie.Title}</span>
                <span className="text-slate-500 text-xs">{formatDate(movie.ReleaseDate)}</span>
            </div>
            
            <div className="flex justify-between text-xs">
                <span className="text-slate-500">{movie.ReleaseStatus}</span>
                <div className="flex">
                    <span className="text-slate-500">{movie.RatingAverage}</span>
                    <AiFillStar className="fill-slate-500 self-center"/>
                </div>
            </div>
        </div>
       
    </li>
  )
}
