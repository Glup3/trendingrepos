import { Selectable } from 'kysely'

export interface Database {
  stars_trend_daily: StarsTrendViewTable
  stars_trend_weekly: StarsTrendViewTable
  stars_trend_monthly: StarsTrendViewTable
}

export interface StarsTrendViewTable {
  github_id: string
  name_with_owner: string
  primary_language: string | null
  description: string | null
  stars_today: number
  stars_prev: number
  stars_diff: number
}

export type StarsTrendView = Selectable<StarsTrendViewTable>
