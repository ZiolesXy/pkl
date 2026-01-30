import React from 'react'

function UserItem() {
  return (
    <div className='flex items-center justify-between gap-2 border rounded-[10px] p-2'>
        <div className="avatar rounded-full min-h-12 min-w-12 bg-amber-300 text-white 
        flex items-center justify-center">
            Ei
        </div>
      <div className='grow'>
        <p className='text-16px font-bold'>EiyuSyaa</p>
        <p className='text-12px text-neutral-400'>EiyuSyaa@gmail.com</p>
      </div>
    </div>
  )
}

export default UserItem
        