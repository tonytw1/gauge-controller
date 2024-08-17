import {Select} from '@headlessui/react'

export function MetricsDropdown({metrics}: { metrics: Metric[] }) {
    return (
        <Select name={"Metric"} aria-description={"Source metric"}>
            {metrics.map(metrics =>
                <option value={metrics.Name}>{metrics.Name} ({metrics.Value})</option>
            )
            }
        </Select>
    )
}
