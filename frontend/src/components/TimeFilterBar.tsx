import { Link, linkOptions } from '@tanstack/react-router'

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

// TODO: fix active link styles when changing page / language
export const TimeFilterBar = () => {
  return (
    <div className="pb-3">
      <div className="flex gap-8 border-b border-[#40484f] px-4">
        {options.map((option) => (
          <Link
            {...option}
            key={option.label}
            className="flex flex-col items-center justify-center border-b-[3px] border-b-transparent pb-[13px] pt-4 text-[#a2abb3]"
            activeProps={{
              className: `text-white border-b-white`,
            }}
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
