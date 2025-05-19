import { useQuery } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { createServerFn } from '@tanstack/react-start'
import { db } from '~/db/database'

export const Route = createFileRoute('/')({
  component: Home,
})

function Home() {
  return (
    <div className="p-2">
      <Component />
    </div>
  )
}

const Component = () => {
  const { data, isPending, isError } = useQuery({
    queryKey: ['stars_trend_monthly'],
    queryFn: () => getStarsTrendMonthly(),
  })

  if (isPending) {
    return <div>Loading...</div>
  }

  if (isError) {
    return <div>Error...</div>
  }

  return (
    <div>
      <table className="table-auto border-separate border-spacing-4">
        <thead>
          <tr>
            <th>Stars</th>
            <th>Name</th>
            <th>Language</th>
            <th>Description</th>
          </tr>
        </thead>

        <tbody>
          {data.map((repo) => (
            <tr key={repo.github_id}>
              <td className="text-right">{repo.diff}</td>
              <td className="line-clamp-1 text-center">
                {repo.name_with_owner}
              </td>
              <td className="text-center">{repo.primary_language}</td>
              <td className="line-clamp-1">{repo.description}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

const getStarsTrendMonthly = createServerFn().handler(async () => {
  return await db
    .selectFrom('stars_trend_monthly')
    .limit(25)
    .selectAll()
    .execute()
})
