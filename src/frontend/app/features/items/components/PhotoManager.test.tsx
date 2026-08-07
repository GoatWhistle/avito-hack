import { fireEvent, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { MAX_PHOTOS } from '#/features/items/lib'
import { PhotoManager } from './PhotoManager'
import { makePhoto, renderWithProviders } from './test-utils'

const listPhotos = vi.fn()
const uploadPhoto = vi.fn()
const deletePhoto = vi.fn()

vi.mock('#/features/items/repository', () => ({
  itemRepository: {
    listPhotos: (id: string) => listPhotos(id),
    uploadPhoto: (params: unknown) => uploadPhoto(params),
    deletePhoto: (itemId: string, photoId: string) =>
      deletePhoto(itemId, photoId),
  },
}))

const makeFile = (name: string, type: string, size: number) => {
  const file = new File(['x'], name, { type })
  Object.defineProperty(file, 'size', { value: size })

  return file
}

beforeEach(() => {
  vi.clearAllMocks()
  listPhotos.mockResolvedValue([])
  uploadPhoto.mockResolvedValue(makePhoto())
  deletePhoto.mockResolvedValue(undefined)
})

describe('PhotoManager', () => {
  it('disables the picker until the item is saved', async () => {
    renderWithProviders(<PhotoManager itemId={undefined} />)

    expect(
      screen.getByText('Сохраните объявление, чтобы добавить фотографии'),
    ).toBeInTheDocument()
    expect(screen.getByLabelText('Добавить фото')).toBeDisabled()
  })

  it('shows the limits in the dropzone hint', async () => {
    renderWithProviders(<PhotoManager itemId="item-1" />)

    expect(
      await screen.findByText('До 10 фото, размер каждого до 5 MB'),
    ).toBeInTheDocument()
    expect(
      screen.getByText(`Осталось ${MAX_PHOTOS} слотов`),
    ).toBeInTheDocument()
  })

  it('uploads accepted files', async () => {
    renderWithProviders(<PhotoManager itemId="item-1" />)

    const input = screen.getByLabelText('Добавить фото')
    await userEvent.upload(
      input,
      makeFile('bike.jpg', 'image/jpeg', 1024 * 1024),
    )

    await waitFor(() => {
      expect(uploadPhoto).toHaveBeenCalledWith(
        expect.objectContaining({ itemId: 'item-1' }),
      )
    })
  })

  it('rejects an unsupported format with a clear message', async () => {
    renderWithProviders(<PhotoManager itemId="item-1" />)

    fireEvent.change(screen.getByLabelText('Добавить фото'), {
      target: { files: [makeFile('doc.pdf', 'application/pdf', 1024)] },
    })

    expect(await screen.findByRole('alert')).toHaveTextContent(
      'doc.pdf: неподдерживаемый формат файла',
    )
    expect(uploadPhoto).not.toHaveBeenCalled()
  })

  it('rejects a file larger than the limit', async () => {
    renderWithProviders(<PhotoManager itemId="item-1" />)

    await userEvent.upload(
      screen.getByLabelText('Добавить фото'),
      makeFile('huge.png', 'image/png', 6 * 1024 * 1024),
    )

    expect(await screen.findByRole('alert')).toHaveTextContent(
      'huge.png: файл больше 5 MB',
    )
    expect(uploadPhoto).not.toHaveBeenCalled()
  })

  it('blocks uploads beyond the photo limit', async () => {
    listPhotos.mockResolvedValue(
      Array.from({ length: MAX_PHOTOS }, (_, index) =>
        makePhoto({ id: `photo-${index}`, position: index }),
      ),
    )
    renderWithProviders(<PhotoManager itemId="item-1" />)

    expect(await screen.findByText('Осталось 0 слотов')).toBeInTheDocument()
    expect(screen.getByLabelText('Добавить фото')).toBeDisabled()
  })

  it('deletes an uploaded photo', async () => {
    listPhotos.mockResolvedValue([makePhoto()])
    renderWithProviders(<PhotoManager itemId="item-1" />)

    const remove = await screen.findByRole('button', { name: 'Удалить фото' })
    await userEvent.click(remove)

    await waitFor(() => {
      expect(deletePhoto).toHaveBeenCalledWith('item-1', 'photo-1')
    })
  })
})
