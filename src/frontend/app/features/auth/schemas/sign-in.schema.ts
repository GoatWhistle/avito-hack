import { z } from 'zod'
import { credentialsShape } from './credentials.schema'

export const SignInSchema = z.object({ ...credentialsShape })
