import Link from "next/link"
import CardTableAxios from "@/components/CardTableAxios"
import { Card } from "@/components/ui/card"
export default async function Page() {

  return (
    <>
    <div>
      <Card>
        <p>Lorem ipsum dolor, sit amet consectetur adipisicing elit. Omnis, blanditiis dolor repudiandae fugit minus dolorem sapiente dolore assumenda reiciendis totam perspiciatis neque praesentium aliquam illum illo voluptate tempora voluptatem asperiores.</p>
      </Card>
    </div>
    <div>
      <Link href="/">Fetch With Native</Link>
      <CardTableAxios />
    </div>

    </>
  )
}
