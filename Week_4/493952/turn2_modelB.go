package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/render"
	"github.com/go-echarts/go-echarts/v2/types"
)

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>E-commerce Metrics Dashboard</title>
			<script src="https://cdnjs.cloudflare.com/ajax/libs/echarts/5.0.2/echarts-en.min.js"></script>
			<style>
				#metrics-chart {
					height: 400px;
				}
			</style>
		</head>
		<body>
			<h1>E-commerce Metrics Dashboard</h1>
			<div id="metrics-chart"></div>
			<script>
				var myChart = echarts.init(document.getElementById('metrics-chart'));
				var xAxisData = [];
				var seriesData = [];

				function fetchMetricsAndUpdateChart() {
					fetch('/metrics/json')
						.then(response => response.json())
						.then(data => {
							// Update xAxisData and seriesData with the latest metrics
							const now = new Date();
							xAxisData.push(now.toLocaleTimeString());
							seriesData.push([data.ConversionRate, data.PageViews, data.BounceRate, data.SessionDuration, data.ClickThroughRate]);
							
							// Limit the data points to keep the chart from getting too wide
							if (seriesData.length > 10) {
								seriesData.shift();
								xAxisData.shift();
							}
							
							myChart.setOption({
								xAxis: {
									type: 'category',
									data: xAxisData