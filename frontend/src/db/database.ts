import { Kysely, PostgresDialect } from 'kysely'
import { Pool } from 'pg'

import { Database } from './types'

const dialect = new PostgresDialect({ pool: new Pool() })

export const db = new Kysely<Database>({ dialect })
