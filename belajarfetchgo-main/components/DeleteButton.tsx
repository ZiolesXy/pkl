"use client"

import { deleteResource } from "@/lib/api/delete"
import { useTransition } from "react"
import { useRouter } from "next/navigation"
import { Button } from "./ui/button"

interface DeleteButtonProps {
  type: "user" | "barang" | "role" | "ownership"
  userId?: number
  barangId?: number
  roleId?: number
  label?: string
  onSuccess?: () => void
}

export default function DeleteButton(props: DeleteButtonProps) {
  const [isPending, startTransition] = useTransition()
  const router = useRouter()

  const handleDelete = () => {
    if (!confirm("Yakin ingin menghapus data ini?")) return

    startTransition(async () => {
      try {
        await deleteResource(props)
        props.onSuccess?.()
        router.refresh()
      } catch (err) {
        const message = err instanceof Error ? err.message : "Delete gagal"
        alert(message)
      }
    })
  }

  return (
    <Button
      onClick={handleDelete}
      disabled={isPending}
      style={{ color: "red" }}
    >
      {props.label ?? "Delete"}
    </Button>
  )
}