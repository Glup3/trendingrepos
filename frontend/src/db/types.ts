import { Selectable } from 'kysely'

export interface Database {
  stars_trend_daily: StarsTrendViewTable
  stars_trend_weekly: StarsTrendViewTable
  stars_trend_monthly: StarsTrendViewTable
  stars: StarsTable
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

export interface StarsTable {
  github_id: string
  stars: number
  time: Date
}

export type StarsTrendView = Selectable<StarsTrendViewTable>
