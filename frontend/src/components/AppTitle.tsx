import { Typography } from '@mui/material'
import { useAtom } from 'jotai'
import { useEffect } from 'react'
import { appTitleState } from '../atoms/settings'
import useFetch from '../hooks/useFetch'

export const AppTitle: React.FC = () => {
  const [appTitle, setAppTitle] = useAtom(appTitleState)

  const { data } = useFetch<{ title: string }>('/webconfig')

  useEffect(() => {
    if (data?.title) {
      setAppTitle(
        data.title.startsWith('"')
          ? data.title.substring(1, data.title.length - 1)
          : data.title
      )
    }
  }, [data])

  return (
    <Typography
      component="h1"
      variant="h6"
      color="inherit"
      noWrap
      sx={{ flexGrow: 1 }}
    >
      {appTitle.startsWith('"') ? appTitle.substring(1, appTitle.length - 1) : appTitle}
    </Typography>
  )
}