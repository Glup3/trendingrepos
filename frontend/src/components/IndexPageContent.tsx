import { Link, getRouteApi } from '@tanstack/react-router'

import { RepoView } from './RepoView'
import { TimeFilterBar } from './TimeFilterBar'

const route = getRouteApi('/')

export const IndexPageContent = () => {
  const search = route.useSearch()

  return (
    <main className="mx-auto flex max-w-screen-xl flex-1 flex-col p-4">
      <div className="mb-4 flex flex-wrap justify-between gap-3">
        <div className="flex min-w-72 flex-col gap-3">
          <p className="tracking-light text-[32px] font-bold leading-tight">
            Trending Repositories
          </p>
          <p className="text-sm font-normal leading-normal">
            Explore the most popular repositories on GitHub, ranked by star
            difference over selectable time periods.
          </p>
        </div>
      </div>

      <TimeFilterBar time={search.time} />

      {search.time === 'monthly' && (
        <div
          className="bg-muted-background mb-4 border p-4 text-sm text-[#6b85ff]"
          role="alert"
        >
          <span className="font-medium">!!!</span>
          <span>
            {' '}
            Monthly data is incomplete, since the crawler has only been running
            since mid of May 2025.
          </span>
        </div>
      )}

      <RepoView />

      <div className="mt-4 flex justify-center gap-4">
        <Link
          from="/"
          search={(prev) => ({
            ...prev,
            page: prev.page - 1,
          })}
          onClick={() => window.scrollTo({ top: 0, behavior: 'smooth' })}
        >
          Prev
        </Link>

        <span className="font-medium">{search.page}</span>

        <Link
          from="/"
          search={(prev) => ({
            ...prev,
            page: prev.page + 1,
          })}
          onClick={() => window.scrollTo({ top: 0, behavior: 'smooth' })}
        >
          Next
        </Link>
      </div>
    </main>
  )
}
