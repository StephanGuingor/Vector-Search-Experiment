import type { Movie } from '@/components/movie'
import type { NextApiRequest, NextApiResponse } from 'next'

import { Client } from '@elastic/elasticsearch'

const client = new Client({
  node: 'https://localhost:9200',
  auth: {
    username: 'elastic',
    password: 'elastic'
  },
  tls: {
    rejectUnauthorized: false,
  }
})

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse<Movie[]>
) {
  const { q } = req.query

  const inferResponse = await client.ml.inferTrainedModel({
    model_id: 'sentence-transformers__msmarco-minilm-l-12-v3',
    body: {
      docs: [
        {
          text_field: q
        },
      ]
    }
  })

  if (!inferResponse.inference_results?.length) {
    res.status(200).json([])
    return
  }
  
  const embedding = inferResponse.inference_results[0].predicted_value as number[]

  const response = await client.search({
    index: 'tmdb-with-embeddings',
    body: {
      knn: {
        "field": "Embedding_Overview.predicted_value",
        "query_vector": embedding,
        "k": 10,
        num_candidates: 100
      }
    }
  })

  const movies = response.hits.hits.map(hit => {
    const source = hit._source as any
    return {
      Title: source.Title,
      TMBdID: source.TMDb_Id,
      ReleaseDate: source.Release_Date,
      RatingAverage: source.Rating_Average,
      RatingCount: source.Rating_Count,
      ReleaseStatus: source.Release_Status,
      Overview: source.Overview
    }
  })

  res.status(200).json(movies)
}
