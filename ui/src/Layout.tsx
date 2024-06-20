import { Menu } from 'antd'
import { DashboardOutlined } from '@ant-design/icons'
import type { MenuProps } from 'antd'
import { useState } from 'react'

type MenuItem = Required<MenuProps>['items'][number]

const loggedInItems: MenuItem[] = [
  {
    label: "PORTal",
    key: "portal"
  },
  {
    type: "divider",
  },
  {
    label: "Dashboard",
    key: "Dashboard",
    icon: <DashboardOutlined />,
  }
]

function Layout() {
  const [current, setCurrent] = useState("portal")
  return (
    <>
  <Menu mode="horizontal" selectedKeys={[current]} items={loggedInItems} />
  {props.children}
  </>
  )
}

export default Home
