import { useQuery } from '@tanstack/react-query'
import { useNavigate } from '@tanstack/react-router'

import { lastUpdateQueryOptions } from '~/utils/api'
import { languages } from '~/utils/searchSchemas'

type Lang = (typeof languages)[number] | null

export const OtherFilters = (props: { language: Lang }) => {
  const navigate = useNavigate({ from: '/' })
  const today = new Date().toLocaleDateString()

  const { data } = useQuery(lastUpdateQueryOptions())

  return (
    <div className="mb-4 flex flex-wrap items-baseline justify-between gap-4">
      <div className="relative">
        <select
          className="bg-muted-background text-foreground border-border appearance-none rounded border py-2 pl-3 pr-8 focus:outline-none focus:ring-2 focus:ring-blue-400 focus:ring-blue-500"
          value={props.language || ''}
          onChange={(e) => {
            navigate({
              from: '/',
              to: '/',
              search: (prev) => ({
                ...prev,
                page: 1,
                language: e.target.value as Lang,
              }),
            })
          }}
        >
          <option value="">All Languages</option>
          {languages.map((lang) => (
            <option key={lang} value={lang}>
              {lang}
            </option>
          ))}
        </select>

        <div className="text-foreground pointer-events-none absolute inset-y-0 right-0 flex items-center px-2">
          <svg
            className="h-4 w-4"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            viewBox="0 0 24 24"
          >
            <path d="M19 9l-7 7-7-7" />
          </svg>
        </div>
      </div>

      <div>
        {data ? (
          <span className="text-muted-foreground text-xs">
            last updated{' '}
            {today === data.time.toLocaleDateString()
              ? data.time.toLocaleTimeString()
              : data.time.toLocaleString()}
          </span>
        ) : null}
      </div>
    </div>
  )
}
