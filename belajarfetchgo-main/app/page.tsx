import Link from "next/link"
import CardTable from "@/components/CardTable"
import { Suspense } from "react"
import { Card } from "@/components/ui/card"
export default async function UsersPage() {

  return (
    <>
      <div>
      <Card>
        <p>Lorem ipsum dolor, sit amet consectetur adipisicing elit. Omnis, blanditiis dolor repudiandae fugit minus dolorem sapiente dolore assumenda reiciendis totam perspiciatis neque praesentium aliquam illum illo voluptate tempora voluptatem asperiores.</p>
      </Card>
    </div>
    <div>
    <Link href="/users-axios">Fetch With Axios</Link>
    <Suspense fallback={<p>Loading...</p>}>
    <CardTable />
    </Suspense>
    </div>
    </>
  )
}
