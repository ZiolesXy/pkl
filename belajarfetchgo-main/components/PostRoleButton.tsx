"use client"

import { Button } from "@/components/ui/button"
import FormDialog from "@/components/FormDialog"
import { postJson } from "@/lib/api/post"

export function PostRoleButton() {
    return (
        <FormDialog<{ name: string }>
            triggerLabel="Post Role"
            title="Create Role"
            fields={[{ name: "name", label: "Role Name", required: true }]}
            onSubmit={async (values) => {
                await postJson("/roles", values)
            }}
        />)
}
