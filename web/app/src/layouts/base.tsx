import { Outlet } from 'react-router-dom'

export default function BaseLayout() {
  return (
    <div className="h-max flex flex-col overflow-hidden">
      <Outlet />
    </div>
  )
}
