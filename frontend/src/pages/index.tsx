import { Search } from '@/components/search'
import { Movie, MovieItem } from '@/components/movie'
import { useState } from 'react'

export default function Home() {
  function onSearch(q: string) {
    fetch(`http://localhost:3000/api/search?q=${q}`)
      .then(response => response.json())
      .then(data => setMovies(data))
  }

  const [movies, setMovies] = useState<Movie[]>([])

  return (
    <main className="flex min-h-screen flex-col p-12 mx-auto">
      <div className='flex self-center w-full max-w-md'>
        <Search onSearch={onSearch} placeholder='Search'/>
      </div>

      <ul className='flex flex-col self-center w-full max-w-md'>
        {movies.map(movie => (
          <MovieItem key={movie.TMBdID} movie={movie} id={movie.TMBdID}/>
        ))}
      </ul>
     
    </main>
  )
}
