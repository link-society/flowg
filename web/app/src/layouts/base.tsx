import { Outlet } from 'react-router-dom'

export const BaseLayout = () => {
  return (
    <div className="h-full flex flex-col overflow-hidden">
      <Outlet />
    </div>
  )
}
