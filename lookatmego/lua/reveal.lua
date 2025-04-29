return {
	plugin = function(input)
		local block = require "lookatmego".block
		local text = input.text

		local idxs = {}
		local chars = {}
		for i = 1, #text do
			local char = text:sub(i, i)
			if char == "\n" then
				table.insert(chars, char)
			elseif char == " " then
				table.insert(chars, " ")
			elseif char == "~" then
				table.insert(chars, " ")
			else
				table.insert(idxs, i)
				table.insert(chars, " ")
			end
		end

		for i = #idxs, 1, -1 do
			local r = math.random(i)
			idxs[i], idxs[r] = idxs[r], idxs[i]
		end

		local out
		for _, idx in ipairs(idxs) do
			local real = text:sub(idx, idx)
			chars[idx] = real

			out = table.concat(chars)
			ch:send(block(out) .. "\n")

			sleep(input.slide_length)
		end

		return block(out) .. "\n"
	end,
}
