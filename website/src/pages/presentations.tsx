import Layout from '@theme/Layout'
import Heading from '@theme/Heading'

const presentations = [
  {
    title: "OSMC",
    date: "November 20th, 2025",
    href: "/flowg/slides/2025-11-20_OSMC/index.html",
  },
]


export default function Presentations(): JSX.Element {
  return (
    <Layout
      title="Presentations"
      description="Index of presentations about FlowG"
    >
      <main>
        <section
          className="container"
          style={{ marginTop: '2rem', marginBottom: '2rem' }}
        >
          <Heading as="h2">Presentations</Heading>

          <ul>
            {presentations.map((presentation) => (
              <li key={presentation.href}>
                <a href={presentation.href}>
                  {presentation.title} <em>({presentation.date})</em>
                </a>
              </li>
            ))}
          </ul>
        </section>
      </main>
    </Layout>
  )
}
