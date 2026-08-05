export const petKeys = {
  all: ['pet'] as const,
  profile: () => [...petKeys.all, 'profile'] as const,
  badges: () => [...petKeys.all, 'badges'] as const,
};
