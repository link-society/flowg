document.addEventListener('htmx:beforeSwap', () => {
  if (window.chart !== undefined) {
    window.chart.destroy()
    window.chart = undefined
  }
})

document.addEventListener('htmx:load', () => {
  const histogramElt  = document.getElementById('data_stream_histogram')
  const histogramData = JSON.parse(histogramElt.dataset.timeserie)

  if (window.chart !== undefined) {
    window.chart.destroy()
  }

  window.chart = new ApexCharts(histogramElt, {
    series: [
      {
        name: 'Logs',
        data: histogramData,
      },
    ],
    chart: {
      type: 'bar',
      width: '100%',
      height: 150,
      animations: {
        enabled: false,
      }
    },
    dataLabels: {
      enabled: false,
    },
    xaxis: {
      type: 'datetime',
    },
  })
  window.chart.render()
})
