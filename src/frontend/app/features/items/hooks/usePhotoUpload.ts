import { useCallback, useState } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import { itemRepository } from '#/features/items/repository'
import { selectPhotos, type PhotoRejection } from '#/features/items/lib'
import { itemKeys } from './query-keys'

export interface UploadTask {
  id: string
  name: string
  progress: number
  failed: boolean
}

const taskId = () =>
  `${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 8)}`

export const usePhotoUpload = (itemId: string | undefined) => {
  const queryClient = useQueryClient()
  const [tasks, setTasks] = useState<UploadTask[]>([])
  const [rejected, setRejected] = useState<PhotoRejection[]>([])

  const patchTask = useCallback((id: string, patch: Partial<UploadTask>) => {
    setTasks((current) =>
      current.map((task) => (task.id === id ? { ...task, ...patch } : task)),
    )
  }, [])

  const dropTask = useCallback((id: string) => {
    setTasks((current) => current.filter((task) => task.id !== id))
  }, [])

  const upload = useCallback(
    async (files: readonly File[], uploadedCount: number) => {
      if (!itemId) return

      const { accepted, rejected: skipped } = selectPhotos(files, uploadedCount)
      setRejected(skipped)
      if (accepted.length === 0) return

      const queued = accepted.map((file) => ({
        file,
        task: { id: taskId(), name: file.name, progress: 0, failed: false },
      }))

      setTasks((current) => [...current, ...queued.map((entry) => entry.task)])

      for (const { file, task } of queued) {
        try {
          await itemRepository.uploadPhoto({
            itemId,
            file,
            onProgress: (progress) => patchTask(task.id, { progress }),
          })
          dropTask(task.id)
        } catch {
          patchTask(task.id, { failed: true })
        }
      }

      await queryClient.invalidateQueries({ queryKey: itemKeys.photos(itemId) })
    },
    [dropTask, itemId, patchTask, queryClient],
  )

  const remove = useCallback(
    async (photoId: string) => {
      if (!itemId) return
      await itemRepository.deletePhoto(itemId, photoId)
      await queryClient.invalidateQueries({ queryKey: itemKeys.photos(itemId) })
    },
    [itemId, queryClient],
  )

  const dismissRejected = useCallback(() => setRejected([]), [])

  return { tasks, rejected, upload, remove, dismissRejected, dropTask }
}
