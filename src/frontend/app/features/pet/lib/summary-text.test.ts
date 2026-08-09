import { afterEach, beforeEach, describe, expect, it } from 'vitest'
import { changeLocale, i18n, initI18n } from '#/i18n'
import type { Advice, SummaryFacts } from '#/features/pet/types'
import { adviceText, summaryMessage } from './summary-text'

const facts = (overrides: Partial<SummaryFacts> = {}): SummaryFacts => ({
  total_xp: 20,
  actions: [{ action: 'item_viewed', count: 3, amount: 20 }],
  level: 4,
  previous_level: 3,
  leveled_up: true,
  xp: 15,
  next_level_xp: 22,
  xp_to_next_level: 7,
  rewards: [],
  badges: [],
  stage: 'baby',
  state: 'neutral',
  satiety: 70,
  happiness: 80,
  energy: 90,
  streak_days: 4,
  streak_broken: false,
  issues_count: 0,
  ...overrides,
})

const t = () => i18n.getFixedT(null, 'catalog')

beforeEach(() => {
  initI18n()
})

afterEach(async () => {
  await changeLocale('ru')
})

describe('summaryMessage', () => {
  it('builds a Russian message from the facts', async () => {
    await changeLocale('ru')
    const message = summaryMessage(t(), facts())

    expect(message).toContain('Ты посмотрел 3 объявления')
    expect(message).toContain('20 XP')
    expect(message).toContain('4 уровня')
  })

  it('builds an English message from the same facts', async () => {
    await changeLocale('en')
    const message = summaryMessage(t(), facts())

    expect(message).toContain('You viewed 3 listings')
    expect(message).toContain('20 XP')
    expect(message).toContain('level 4')
    expect(message).not.toMatch(/[Ѐ-ӿ]/)
  })

  it('uses the empty day wording when nothing happened', async () => {
    await changeLocale('en')
    const message = summaryMessage(
      t(),
      facts({ total_xp: 0, actions: [], leveled_up: false }),
    )

    expect(message).toContain('I missed you all day')
    expect(message).toContain("we'll catch up together")
  })

  it('reports the leaderboard rank when known', async () => {
    await changeLocale('en')
    const message = summaryMessage(t(), facts({ leaderboard_rank: 7 }))

    expect(message).toContain('holding 7')
  })
})

describe('adviceText', () => {
  const advice: Advice = {
    text: 'У объявления «Велосипед» нет ни одной фотографии',
    item_id: 'item-1',
    item_title: 'Велосипед',
    action: 'add_photo',
  }

  it('translates the advice into English with the listing title', async () => {
    await changeLocale('en')

    expect(adviceText(t(), advice)).toBe(
      '"Велосипед" has no photos at all — listings with photos sell noticeably faster.',
    )
  })

  it('uses the generic wording when there is no listing title', async () => {
    await changeLocale('en')

    const text = adviceText(t(), { ...advice, item_title: null })

    expect(text).toBe(
      'One of your listings has no photos at all — listings with photos sell noticeably faster.',
    )
    expect(text).not.toContain('""')
  })

  it('translates the check-in advice that never carries a title', async () => {
    await changeLocale('en')

    const text = adviceText(t(), {
      text: 'Объявления у тебя в порядке',
      item_id: null,
      item_title: null,
      action: 'check_in',
    })

    expect(text).toContain('Your listings are in good shape')
    expect(text).not.toMatch(/[Ѐ-ӿ]/)
  })

  it('falls back to the server text for an unknown action', async () => {
    await changeLocale('en')

    expect(adviceText(t(), { ...advice, action: 'something_new' })).toBe(
      advice.text,
    )
  })
})
