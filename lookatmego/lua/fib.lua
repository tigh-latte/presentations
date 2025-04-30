return {
	plugin = function()
		local seq = { 0, 1 }

		for _ = 3, 15 do
			table.insert(seq, seq[#seq] + seq[#seq - 1])
		end

		return table.concat(seq, " ")
	end,
}
