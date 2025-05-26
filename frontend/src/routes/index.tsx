import { createFileRoute, stripSearchParams } from '@tanstack/react-router'
import { zodValidator } from '@tanstack/zod-adapter'

import { IndexPageContent } from '~/components/IndexPageContent'
import { defaultPage, indexSearchSchema } from '~/utils/searchSchemas'

export const Route = createFileRoute('/')({
  component: IndexPageContent,
  validateSearch: zodValidator(indexSearchSchema),
  search: {
    middlewares: [stripSearchParams({ page: defaultPage })],
  },
})
