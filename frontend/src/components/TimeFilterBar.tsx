import { Link, linkOptions } from '@tanstack/react-router'
import { twMerge } from 'tailwind-merge'

const options = linkOptions([
  {
    from: '/',
    label: 'Daily',
    search: { page: 1, time: 'daily' },
    activeOptions: { exact: true, includeSearch: true },
  },
  {
    from: '/',
    label: 'Weekly',
    search: { page: 1, time: 'weekly' },
    activeOptions: { exact: true, includeSearch: true },
  },
  {
    from: '/',
    label: 'Monthly',
    search: { page: 1, time: 'monthly' },
    activeOptions: { exact: true, includeSearch: true },
  },
])

export const TimeFilterBar = (props: { time: string }) => {
  return (
    <div className="pb-4">
      <div className="flex gap-8 border-b border-[#40484f]">
        {options.map((option) => (
          <Link
            {...option}
            key={option.label}
            className={twMerge(
              'flex flex-col items-center justify-center border-b-[3px] border-b-transparent pb-[13px] pt-4 text-[#a2abb3] hover:text-white',
              props.time === option.search.time && 'border-b-white text-white',
            )}
          >
            <p className="text-sm font-bold leading-normal tracking-[0.015em]">
              {option.label}
            </p>
          </Link>
        ))}
      </div>
    </div>
  )
}
